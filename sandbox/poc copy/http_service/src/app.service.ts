import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHome(): string {
    return `${process.env.SERVICE_NAME}-${process.env.SERVICE_VERSION}`;
  }

  getFoo(): string {
    return `foo ${process.env.SERVICE_NAME}`;
  }

  getBar(): string {
    return `bar ${process.env.SERVICE_NAME}`;
  }

  getBaz(): string {
    return `baz ${process.env.SERVICE_NAME}`;
  }
}
