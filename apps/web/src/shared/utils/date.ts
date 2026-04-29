export function formatDate(value:string):string{const d=new Date(value);return Number.isNaN(d.getTime())?"Unknown":d.toLocaleDateString();}
