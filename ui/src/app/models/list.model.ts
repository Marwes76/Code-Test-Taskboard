import { Task } from './task.model';

export interface List {
	uuid:		string;
	title:		string;
	description:	string;
	sortOrder:	number;
	createdAt:	string;
	updatedAt:	string;

	// Relations
	tasks:		Task[];
}