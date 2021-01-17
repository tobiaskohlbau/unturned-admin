import { Token } from '../models';

function getCookie(cname: string): string {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for(var i = 0; i <ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

export function isActivated(): boolean | null {
  const token = getToken();
  if (token === null) {
    return null;
  }
  return token.activated;
}

export function isAuthenticated(): boolean {
  if (getCookie("tid") !== "") {
    return true;
  }
  return false;
}

export function getToken(): Token | null {
  const cookie: string = getCookie("tid");
  if (cookie == "") {
    return null;
  }
  return JSON.parse(atob(cookie));
}

export function hasPermission(permission: string): boolean {
  const token = getToken();
  if (token == null) {
    return false;
  }
  if (token.permissions === undefined) {
    return false;
  }
  if (token.permissions.indexOf(permission) === -1) {
    return false;
  }
  return true;
}
