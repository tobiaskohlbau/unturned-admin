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

export function isAuthenticated(): boolean {
  if (getCookie("tid") !== "") {
    return true;
  }
  return false;
}

interface Token {
  username: string
  permissions: string[]
}

export function hasPermission(permission: string): boolean {
  const cookie: string = getCookie("tid");
  if (cookie == "") {
    return false;
  }
  const token: Token = JSON.parse(atob(cookie));
  if (token.permissions === undefined) {
    return false;
  }
  if (token.permissions.indexOf(permission) === -1) {
    return false;
  }
  return true;
}
