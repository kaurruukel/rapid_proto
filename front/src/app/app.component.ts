import { Component } from '@angular/core';
import { HeaderComponent } from './shared/header/header.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [HeaderComponent, DashboardComponent],
    templateUrl: './app.component.html',
    styleUrl: './app.component.less',
})
export class AppComponent {
    title = 'front';
}
