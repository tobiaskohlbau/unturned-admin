export interface File {
  path: string
  name: string
  content_type: string
}

export interface Token {
  username: string
  permissions: string[]
  activated: boolean
}

export interface User {
  username: string
  permissions: string[]
  activated: boolean
}
