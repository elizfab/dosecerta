import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api/api-service';
import { Exam, CreateExamInput } from '../../../shared/types/exam.interface';

/** Consome os endpoints /api/v1/exams do backend. */
@Injectable({ providedIn: 'root' })
export class ExamService {
  private readonly api = inject(ApiService);
  private readonly path = '/api/v1/exams';

  list(): Observable<Exam[]> {
    return this.api.get<Exam[]>(this.path);
  }

  getById(id: string): Observable<Exam> {
    return this.api.get<Exam>(`${this.path}/${id}`);
  }

  create(input: CreateExamInput): Observable<Exam> {
    return this.api.post<Exam>(this.path, input);
  }

  update(id: string, input: Partial<CreateExamInput>): Observable<Exam> {
    return this.api.put<Exam>(`${this.path}/${id}`, input);
  }

  remove(id: string): Observable<void> {
    return this.api.delete<void>(`${this.path}/${id}`);
  }
}
