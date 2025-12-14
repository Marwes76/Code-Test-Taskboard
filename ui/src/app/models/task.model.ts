export class Task {
	uuid:		string = "";
	listUuid:	string = "";
	title:		string = "";
	description:	string = "";
	sortOrder:	number = -1;
	createdAt:	string = "";
	updatedAt:	string = "";

	constructor(init?: Partial<Task>) {
		Object.assign(this, init);
	}
}
