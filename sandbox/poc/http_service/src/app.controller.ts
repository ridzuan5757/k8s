import { Controller, Get } from '@nestjs/common';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getHome(): string {
    return this.appService.getHome();
  }

  @Get('foo')
  getFoo(): string {
    return this.appService.getFoo();
  }

  @Get('bar')
  getBar(): string {
    return this.appService.getBar();
  }

  @Get('baz')
  getBaz(): string {
    return this.appService.getBaz();
  }
}
