import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BehaviorSubject, Observable } from 'rxjs';
import { List } from '../../models/list.model';
import { ListService } from '../../services/list.service';
import { ListComponent } from '../list/list.component';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
	selector: 'app-board',
	standalone: true,
	imports: [
		CommonModule,
		ListComponent,
		MatButtonModule,
		MatIconModule,
	],
	templateUrl: './board.component.html',
	styleUrls: ['./board.component.scss']
})

export class BoardComponent {
	private listsSubject = new BehaviorSubject<List[]>([]);
	lists: Observable<List[]> = this.listsSubject.asObservable();
	lastListIndex: number = 0;
	newList: Partial<List> = {};

	constructor(private listService: ListService) {
		this.listService.getAllLists().subscribe(lists => {
			this.listsSubject.next(lists);
			this.lastListIndex = lists.length - 1;
		});
	}

	onAddList() {
		const list = new List({
			sortOrder: this.listsSubject.value.length
		});
		const lists = [ ...this.listsSubject.value, list ];
		this.listsSubject.next(lists);
		this.lastListIndex = lists.length - 1;
	}

	moveUpListAtIndex(index: number) {
		const lists = [ ...this.listsSubject.value ];
		var list = lists[index];
		var priorList = lists[index - 1];
		const tempSortOrder = list.sortOrder;
		list.sortOrder = priorList.sortOrder;
		priorList.sortOrder = tempSortOrder;

		const listsToUpdate: { [uuid: string]: List } = {
			[list.uuid] : list,
			[priorList.uuid] : priorList,
		};
		this.listService.updateLists(listsToUpdate).subscribe({
			next: (listsUpdated: { [uuid: string]: List }) => {
				list = listsUpdated[list.uuid];
				priorList = listsUpdated[priorList.uuid];
				lists[index - 1] = list;
				lists[index] = priorList;
				this.listsSubject.next(lists);
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to update lists', err);
			}
		});
	}

	moveDownListAtIndex(index: number) {
		const lists = [ ...this.listsSubject.value ];
		var list = lists[index];
		var nextList = lists[index + 1];
		const tempSortOrder = list.sortOrder;
		list.sortOrder = nextList.sortOrder;
		nextList.sortOrder = tempSortOrder;

		const listsToUpdate: { [uuid: string]: List } = {
			[list.uuid] : list,
			[nextList.uuid] : nextList,
		};
		this.listService.updateLists(listsToUpdate).subscribe({
			next: (listsUpdated: { [uuid: string]: List }) => {
				list = listsUpdated[list.uuid];
				nextList = listsUpdated[nextList.uuid];
				lists[index + 1] = list;
				lists[index] = nextList;
				this.listsSubject.next(lists);
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to update lists', err);
			}
		});
	}

	updateListAtIndex(index: number, list: List) {
		const lists = [ ...this.listsSubject.value ];
		lists[index] = list;
		this.listsSubject.next(lists);
	}

	deleteListAtIndex(index: number) {
		const lists = [ ...this.listsSubject.value ];
		lists.splice(index, 1);
		this.listsSubject.next(lists);
		this.lastListIndex = lists.length - 1;
	}
}
