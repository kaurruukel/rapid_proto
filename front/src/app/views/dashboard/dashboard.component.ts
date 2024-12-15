import { Component, computed } from '@angular/core';
import { NgForOf, NgIf } from '@angular/common';
import { MailEventsService } from '../../mail-events.service';
import { MailEventComponent } from '../../shared/mail-event/mail-event.component';

@Component({
    selector: 'app-dashboard',
    standalone: true,
    imports: [NgIf, NgForOf, MailEventComponent],
    templateUrl: './dashboard.component.html',
    styleUrl: './dashboard.component.less',
    providers: [MailEventsService]
})
export class DashboardComponent {

    unSeen = computed<number>(() => {
        return this.mailEventsService.mailEvents()?.filter(m => !m.acknowledged).length || 0;
    });


    constructor(public mailEventsService: MailEventsService) {
        this.mailEventsService.getEvents()
        setInterval(() => this.mailEventsService.getEvents(),
            2000
        )
    }
}
