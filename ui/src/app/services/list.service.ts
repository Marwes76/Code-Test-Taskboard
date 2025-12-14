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

	getList(uuid: string): Observable<List> {
		return this.httpService.get<List>(`lists/${uuid}`);
	}

	createList(list: Partial<List>) {
		return this.httpService.post<List>('lists', list);
	}

	updateLists(lists: { [uuid: string]: List }) {
		return this.httpService.put<{ [uuid: string]: List }>(`lists`, lists);
	}

	deleteList(uuid: string) {
		return this.httpService.delete(`lists/${uuid}`);
	}
}
