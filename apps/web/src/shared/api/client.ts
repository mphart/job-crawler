export type HttpMethod="GET"|"POST"|"PATCH";
export async function requestJson<T>(url:string,method:HttpMethod="GET",body?:unknown):Promise<T>{// DECISION: APIs are mocked until backend contracts land.
const response=await fetch(url,{method,headers:{"Content-Type":"application/json"},body:body?JSON.stringify(body):undefined});if(!response.ok)throw new Error(`Request failed: ${response.status}`);return await response.json() as T;}
