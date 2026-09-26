import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { environment } from '../../../../app/environment/environment';

import { HealthService } from './health-service';

describe('HealthService', () => {
  let service: HealthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(HealthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should GET /health', () => {
    service.check().subscribe((res) => {
      expect(res.status).toBe('ok');
    });

    const req = httpMock.expectOne(`${environment.apiUrl}/health`);
    expect(req.request.method).toBe('GET');
    req.flush({ status: 'ok', service: 'backend' });
  });
});
