import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Task } from '../../models/task.model';
import { TaskService } from '../../services/task.service';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { EditState } from '../../helpers/edit-state.class';
import { HttpErrorResponse } from '@angular/common/http';
import { FormsModule } from '@angular/forms';


@Component({
	selector: 'app-task',
	standalone: true,
	imports: [
		CommonModule,
		FormsModule,
		MatButtonModule,
		MatFormFieldModule,
		MatIconModule,
		MatInputModule,
	],
	templateUrl: './task.component.html',
	styleUrls: ['./task.component.scss']
})
export class TaskComponent {
	@Input() index!: number;
	@Input() task!: Task;
	@Input() listUuid!: string;
	@Input() isFirstList!: boolean;
	@Input() isLastList!: boolean;
	@Output() moveUp = new EventEmitter<number>();
	@Output() moveDown = new EventEmitter<number>();
	@Output() moveLeft = new EventEmitter<number>();
	@Output() moveRight = new EventEmitter<number>();
	@Output() update = new EventEmitter<{index: number; task: Task}>();
	@Output() delete = new EventEmitter<number>();

	editableTask!: Task;
	editState: EditState = EditState.DEFAULT;

	constructor(private taskService: TaskService) {}

	ngOnInit() {
		this.editableTask = new Task({ ...this.task });
		if (this.editableTask.uuid === "") {
			this.editState = EditState.NEW;
			this.editableTask.listUuid = this.listUuid;
		}
	}

	onEdit() {
		this.editState = EditState.EDITING;
	}

	onMoveLeft() {
		this.moveLeft.emit(this.index);
	}

	onMoveRight() {
		this.moveRight.emit(this.index);
	}

	onSave() {
		if (this.editState.isNew()) {
			this.editState = EditState.SAVING;

			this.taskService.createTask(this.editableTask).subscribe({
				next: (task: Task) => {
					this.update.emit({index: this.index, task: task});
					this.editState = EditState.DEFAULT;
				},
				error: (err: HttpErrorResponse) => {
					console.error('Failed to create task', err);
					this.editState = EditState.NEW;
				}
			});

			this.editState = EditState.NEW;
		} else if (this.editState.isEditing()) {
			this.editState = EditState.SAVING;

			const tasks: { [uuid: string]: Task } = { [this.editableTask.uuid] : this.editableTask };
			this.taskService.updateTasks(tasks).subscribe({
				next: (tasks: { [uuid: string]: Task }) => {
					const task = tasks[this.editableTask.uuid];
					if (task !== undefined) {
						this.update.emit({index: this.index, task: task});
						this.editState = EditState.DEFAULT;
					} else {
						console.error('Failed to update task');
						this.editState = EditState.DEFAULT;
					}
				},
				error: (err: HttpErrorResponse) => {
					console.error('Failed to update task', err);
					this.editState = EditState.DEFAULT;
				}
			});
		}
	}

	onCancel() {
		if (this.editState.isNew()) {
			this.delete.emit(this.index);
		} else if (this.editState.isEditing()) {
			this.editableTask = { ...this.task };
			this.editState = EditState.DEFAULT;
		}
	}

	onDelete() {
		this.editState = EditState.SAVING;

		this.taskService.deleteTask(this.task.uuid).subscribe({
			next: () => {
				this.delete.emit(this.index);
			},
			error: (err: HttpErrorResponse) => {
				console.error('Failed to delete task', err);
			}
		});

		this.editState = EditState.DEFAULT;
	}
}
