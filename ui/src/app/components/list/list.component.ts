import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { List } from '../../models/list.model';
import { TaskComponent } from '../task/task.component';

@Component({
	selector: 'app-list',
	standalone: true,
	imports: [
		CommonModule,
		TaskComponent
	],
	templateUrl: './list.component.html',
	styleUrls: ['./list.component.scss']
})

export class ListComponent {
	@Input() list!: List;
}
