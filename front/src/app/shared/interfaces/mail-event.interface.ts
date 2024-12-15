
export interface MailEvent {
    id: string;
    acknowledged: boolean;
    date: string;
    type: 'inserted' | 'removed'
}
