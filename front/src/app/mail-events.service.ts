import { Injectable, signal } from '@angular/core';
import { Observable, of } from 'rxjs';
import { MailEvent } from './shared/interfaces/mail-event.interface';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';


@Injectable({
    providedIn: 'root'
})
export class MailEventsService {
    constructor(
        private http: HttpClient
    ) { }
    mailEvents = signal<MailEvent[] | null>(null);

    public getEvents(): void {
        this.http.get<MailEvent[]>(`${environment.backendUrl}/events`)
            .subscribe({
                next: (events) => {
                    console.log(events);
                    this.mailEvents.set(events);
                },
                error: (err) => {
                    console.log(err);
                }
            })

    }

    public ackEvent(id: string): void {
        this.http.patch(`${environment.backendUrl}/events/${id}`, {
            acknowledged: true,
        })
            .subscribe({
                next: () => {
                    console.log("Event acknowledged")
                },
                error: (err) => {
                    console.log(err);
                }
            })
    }
}
