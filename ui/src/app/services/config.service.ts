import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom, of } from 'rxjs';
import { catchError } from 'rxjs/operators';

interface AppConfig {
	apiUrl?:	string;
}

@Injectable({ providedIn: 'root' })
export class ConfigService {
	private readonly http = inject(HttpClient);
	private _config!: AppConfig;

	async loadConfig(): Promise<void> {
		const config = await firstValueFrom(
			this.http.get<AppConfig>('/assets/config.json').pipe(
				catchError(err => {
					console.error('Failed to load config', err);
					return of({} as AppConfig);
				})
			)
  		);
		this._config = config;
	}

	get apiUrl(): string {
		if (!this._config) {
			throw new Error('Config not loaded yet');
		}
		if (!this._config.apiUrl) {
			throw new Error('apiUrl not set in config');
		}
		return this._config.apiUrl;
	}
}
