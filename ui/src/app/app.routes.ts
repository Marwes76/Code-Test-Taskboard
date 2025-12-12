import { Routes } from '@angular/router';
import { BoardComponent } from './components/board/board.component';

export const routes: Routes = [
	{
		// Redirect root to board for now, since there is nothing else
		path: '',
		redirectTo: 'board',
		pathMatch: 'full'
	},
	{
		path: 'board',
		component: BoardComponent
	}
];
