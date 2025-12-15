import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from './config.service';

@Injectable({ providedIn: 'root' })
export class HttpService {

	private readonly http = inject(HttpClient);
	private readonly config = inject(ConfigService);

	constructor() {}

	get<T>(endpoint: string, params: { params?: HttpParams, headers?: HttpHeaders } = {}): Observable<T> {
		return this.http.get<T>(`${this.config.apiUrl}/${endpoint}`, params);
	}

	post<T>(endpoint: string, body: any, headers?: HttpHeaders): Observable<T> {
		return this.http.post<T>(`${this.config.apiUrl}/${endpoint}`, body, { headers });
	}

	put<T>(endpoint: string, body: any, headers?: HttpHeaders): Observable<T> {
		return this.http.put<T>(`${this.config.apiUrl}/${endpoint}`, body, { headers });
	}

	delete<T>(endpoint: string, headers?: HttpHeaders): Observable<T> {
		return this.http.delete<T>(`${this.config.apiUrl}/${endpoint}`, { headers });
	}

	getHttpParams(params: { [key: string]: string }): HttpParams {
		let httpParams = new HttpParams();

		for (const key in params) {
			if (params[key] && params[key] !== "") {
				httpParams = httpParams.set(key, params[key]);
			}
		}

		return httpParams;
	}
}
