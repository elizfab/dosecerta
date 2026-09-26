// Setup do jest-preset-angular v17 (inicializa o TestBed do Angular para os testes).
// App é zoneless (sem zone.js), por isso o preset "zoneless".
import { setupZonelessTestEnv } from 'jest-preset-angular/setup-env/zoneless';

setupZonelessTestEnv();
