import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { Component } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { DevelopmentUserHeaderService } from '$core/dev/development-user-header.service';
import { MainLayout } from './main-layout';

@Component({
  selector: 'app-test-page',
  template: '<p class="text-gray-700 dark:text-gray-300">Page content goes here</p>',
})
class TestPage {}

describe('MainLayout screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const providers = [
    provideTranslocoForTest(),
    provideRouter([{ path: '**', component: TestPage }]),
    DevelopmentUserHeaderService,
  ];

  describe('desktop', () => {
    beforeEach(async () => {
      await page.viewport(1280, 720);
    });

    it('light', async () => {
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('desktop-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('desktop-dark');
    });
  });

  describe('mobile', () => {
    beforeEach(async () => {
      await page.viewport(375, 667);
    });

    it('light', async () => {
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('mobile-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('mobile-dark');
    });
  });

  describe('mobile menu open', () => {
    beforeEach(async () => {
      await page.viewport(375, 667);
    });

    it('light', async () => {
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      await page.getByRole('button', { name: /toggleMenu/i }).click();
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('mobile-menu-open-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(MainLayout, { providers });
      await fixture.whenStable();
      await page.getByRole('button', { name: /toggleMenu/i }).click();
      await fixture.whenStable();
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('mobile-menu-open-dark');
    });
  });
});
