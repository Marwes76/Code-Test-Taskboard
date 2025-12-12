import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { List } from '../models/list.model';
import { HttpService } from './http.service';

@Injectable({ providedIn: 'root' })
export class ListService {

	constructor(private httpService: HttpService) {}

	getAllLists(): Observable<List[]> {
		return this.httpService.get<List[]>('lists');
	}
}
