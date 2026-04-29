export function isEmail(value:string):boolean{return /.+@.+\..+/.test(value);}
export function hasMinLength(value:string,min:number):boolean{return value.trim().length>=min;}
