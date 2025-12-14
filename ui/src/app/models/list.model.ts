import { Task } from './task.model';

export class List {
	uuid:		string = "";
	title:		string = "";
	description:	string = "";
	sortOrder:	number = -1;
	createdAt:	string = "";
	updatedAt:	string = "";

	// Relations
	tasks:		Task[] = [];

	constructor(init?: Partial<List>) {
		Object.assign(this, init);
	}
}
