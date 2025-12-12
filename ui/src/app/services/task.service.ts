import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Task } from '../models/task.model';
import { HttpService } from './http.service';

@Injectable({ providedIn: 'root' })
export class TaskService {

	constructor(private httpService: HttpService) {}
}
