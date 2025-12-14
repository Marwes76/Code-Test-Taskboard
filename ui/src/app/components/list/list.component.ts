import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BehaviorSubject, Observable } from 'rxjs';
import { List } from '../../models/list.model';
import { ListService } from '../../services/list.service';
import { Task } from '../../models/task.model';
import { TaskComponent } from '../task/task.component';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { EditState } from '../../helpers/edit-state.class';
import { HttpErrorResponse } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

@Component({
	selector: 'app-list',
	standalone: true,
	imports: [
		CommonModule,
		FormsModule,
		MatButtonModule,
		MatFormFieldModule,
		MatIconModule,
		MatInputModule,
		TaskComponent,
	],
	templateUrl: './list.component.html',
	styleUrls: ['./list.component.scss']
})

export class ListComponent {
	@Input() index!: number;
	@Input() lastIndex!: number;
	@Input() list!: List;
	@Output() moveUp = new EventEmitter<number>();
	@Output() moveDown = new EventEmitter<number>();
	@Output() update = new EventEmitter<{index: number; list: List}>();
	@Output() delete = new EventEmitter<number>();

	editableList!: List;
	editState: EditState = EditState.DEFAULT;

	private tasksSubject = new BehaviorSubject<Task[]>([]);
	tasks: Observable<Task[]> = this.tasksSubject.asObservable();
	lastTaskIndex: number = 0;
	newTask: Partial<Task> = {};

	constructor(private listService: ListService) {}

	ngOnInit() {
		const tasks = this.list.tasks.map(t => new Task(t));
		this.editableList = new List({
			...this.list,
			tasks: tasks
		});
		if (this.editableList.uuid === "") {
			this.editState = EditState.NEW;
		}

		this.tasksSubject.next(tasks);
		this.lastTaskIndex = tasks.length - 1;
	}

	onAddTask() {
		const task = new Task({
			sortOrder: this.tasksSubject.value.length
		});
		const tasks = [ ...this.tasksSubject.value, task ];
		this.tasksSubject.next(tasks);
		this.lastTaskIndex = tasks.length - 1;
	}

	updateTaskAtIndex(index: number, task: Task) {
		const tasks = [ ...this.tasksSubject.value ];
		tasks[index] = task;
		this.tasksSubject.next(tasks);
	}

	deleteTaskAtIndex(index: number) {
		const tasks = [ ...this.tasksSubject.value ];
		tasks.splice(index, 1);
		this.tasksSubject.next(tasks);
		this.lastTaskIndex = tasks.length - 1;
	}

	onEdit() {
		this.editState = EditState.EDITING;
	}

	onMoveUp() {
		this.moveUp.emit(this.index);
	}

	onMoveDown() {
		this.moveDown.emit(this.index);
	}

	onSave() {
		if (this.editState.isNew()) {
			this.editState = EditState.SAVING;

			this.listService.createList(this.editableList).subscribe({
				next: (list: List) => {
					this.update.emit({index: this.index, list: list});
					this.editState = EditState.DEFAULT;
				},
				error: (err: HttpErrorResponse) => {
					console.error('Failed to create list', err);
					this.editState = EditState.NEW;
				}
			});

			this.editState = EditState.NEW;
		} else if (this.editState.isEditing()) {
			this.editState = EditState.SAVING;

			const lists: { [uuid: string]: List } = { [this.editableList.uuid] : this.editableList };
			this.listService.updateLists(lists).subscribe({
				next: (lists: { [uuid: string]: List }) => {
					const list = lists[this.editableList.uuid];
					if (list !== undefined) {
						this.update.emit({index: this.index, list: list});
						this.editState = EditState.DEFAULT;
					} else {
						console.error('Failed to update list');
						this.editState = EditState.DEFAULT;
					}
				},
				error: (err: HttpErrorResponse) => {
					console.error('Failed to update list', err);
					this.editState = EditState.DEFAULT;
				}
			});
		}
	}

	onCancel() {
		if (this.editState.isNew()) {
			this.delete.emit(this.index);
		} else if (this.editState.isEditing()) {
			this.editableList = { ...this.list };
			this.editState = EditState.DEFAULT;
		}
	}

	onDelete() {
		this.editState = EditState.SAVING;

		this.listService.deleteList(this.list.uuid).subscribe({
			next: () => {
				this.delete.emit(this.index);
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to delete list', err);
			}
		});

		this.editState = EditState.DEFAULT;
	}
}
