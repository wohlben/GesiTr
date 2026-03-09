import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  host: { class: 'block' },
  template: '<router-outlet />',
})
export class App {}
