import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { List } from '../../models/list.model';
import { ListService } from '../../services/list.service';
import { ListComponent } from '../list/list.component';

import { Observable } from 'rxjs';

@Component({
	selector: 'app-board',
	standalone: true,
	imports: [
		CommonModule,
		ListComponent,
	],
	templateUrl: './board.component.html',
	styleUrls: ['./board.component.scss']
})

export class BoardComponent {
	lists: Observable<List[]>;

	constructor(private listService: ListService) {
		this.lists = this.listService.getAllLists();
	}
}
