import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BehaviorSubject, Observable } from 'rxjs';
import { List } from '../../models/list.model';
import { ListService } from '../../services/list.service';
import { ListComponent } from '../list/list.component';
import { Task } from '../../models/task.model';
import { TaskService } from '../../services/task.service';
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
	newList: Partial<List> = {};

	constructor(private listService: ListService, private taskService: TaskService) {
		this.listService.getAllLists().subscribe(lists => {
			this.listsSubject.next(lists);
		});
	}

	onAddList() {
		const list = new List({
			sortOrder: this.listsSubject.value.length
		});
		const lists = [ ...this.listsSubject.value, list ];
		this.listsSubject.next(lists);
	}

	moveListUpAtIndex(index: number) {
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

	moveListDownAtIndex(index: number) {
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
	}

	moveTaskLeftInList(index: number, task: Task) {
		const lists = [ ...this.listsSubject.value ];
		var priorList = lists[index - 1];
		task.listUuid = priorList.uuid;

		var priorTasks = priorList.tasks.map(t => new Task(t));
		task.sortOrder = priorTasks.length;

		const tasks: { [uuid: string]: Task } = { [task.uuid] : task};
		this.taskService.updateTasks(tasks).subscribe({
			next: (tasks: { [uuid: string]: Task }) => {
				const updatedTask = tasks[task.uuid];
				if (updatedTask !== undefined) {
					// We keep things simple and just reload all lists for now
					this.listService.getAllLists().subscribe({
						next: (lists: List[]) => {
							this.listsSubject.next(lists);
						},
						error: (err: HttpErrorResponse) => {
							console.error('Failed to get lists', err);
						}
					});
				} else {
					console.error('Failed to update task');
				}
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to update task', err);
			}
		});

	}

	moveTaskRightInList(index: number, task: Task) {
		const lists = [ ...this.listsSubject.value ];
		var nextList = lists[index + 1];
		task.listUuid = nextList.uuid;

		var nextTasks = nextList.tasks.map(t => new Task(t));
		task.sortOrder = nextTasks.length;

		const tasks: { [uuid: string]: Task } = { [task.uuid] : task};
		this.taskService.updateTasks(tasks).subscribe({
			next: (tasks: { [uuid: string]: Task }) => {
				const updatedTask = tasks[task.uuid];
				if (updatedTask !== undefined) {
					// We keep things simple and just reload all lists for now
					this.listService.getAllLists().subscribe({
						next: (lists: List[]) => {
							this.listsSubject.next(lists);
						},
						error: (err: HttpErrorResponse) => {
							console.error('Failed to get lists', err);
						}
					});
				} else {
					console.error('Failed to update tasks');
				}
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to update task', err);
			}
		});
	}
}
