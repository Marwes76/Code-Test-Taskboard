import { bootstrapApplication } from '@angular/platform-browser';
import { App } from './app/app';
import { ConfigService } from './app/services/config.service';
import { appConfig } from './app/app.config';
import { APP_INITIALIZER, importProvidersFrom } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

export function loadConfig(config: ConfigService) {
	return () => config.loadConfig();
}

bootstrapApplication(App, {
	providers: [
		ConfigService,
		...appConfig?.providers ?? [],
		{
			provide: APP_INITIALIZER,
			useFactory: loadConfig,
			deps: [ConfigService],
			multi: true
		},
		importProvidersFrom(HttpClientModule)
	]
}).catch(err => console.error(err));
