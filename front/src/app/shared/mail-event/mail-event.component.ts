import { Component, Input } from '@angular/core';
import { MailEvent } from '../interfaces/mail-event.interface';
import { NgClass } from '@angular/common';
import { MailEventsService } from '../../mail-events.service';
import { DatePipe } from '@angular/common';

@Component({
    selector: 'app-mail-event',
    standalone: true,
    imports: [NgClass, DatePipe],
    templateUrl: './mail-event.component.html',
    styleUrl: './mail-event.component.less',
})
export class MailEventComponent {
    @Input() mailEvent!: MailEvent;

    constructor(private mailEventsService: MailEventsService) { }

    public ackEvent(): void {
        this.mailEventsService
            .ackEvent(this.mailEvent.id)
        this.mailEvent.acknowledged = true;
    }

    public getDescription(type: MailEvent['type']): string {
        switch (type) {
            case 'inserted':
                return 'Mail has been inserted';
            case 'removed':
                return 'Mail has been removed';
        }
    }
}
