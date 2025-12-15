import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Task } from '../models/task.model';
import { HttpService } from './http.service';

@Injectable({ providedIn: 'root' })
export class TaskService {

	constructor(private httpService: HttpService) {}

	searchTasks(params: { [key: string]: string }): Observable<Task[]> {
		const httpParams = this.httpService.getHttpParams(params);
		return this.httpService.get<Task[]>(`tasks/search`, { params: httpParams });
	}

	getTask(uuid: string): Observable<Task> {
		return this.httpService.get<Task>(`tasks/${uuid}`);
	}

	createTask(task: Partial<Task>) {
		return this.httpService.post<Task>('tasks', task);
	}

	updateTasks(tasks: { [uuid: string]: Task }) {
		return this.httpService.put<{ [uuid: string]: Task }>(`tasks`, tasks);
	}

	deleteTask(uuid: string) {
		return this.httpService.delete(`tasks/${uuid}`);
	}
}
