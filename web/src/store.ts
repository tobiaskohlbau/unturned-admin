import { inject, onMounted, reactive, readonly } from "vue";
import { getToken } from './utils';

import { Token } from "./models";

export const storeSymbol = Symbol('store');

export const createStore = () => {
  const state = reactive({ token: {} });
  const setToken = (token: Token) => state.token = token;

  return { setToken, state: readonly(state) };
}

interface State {
  token: Token
}

interface Store {
  setToken(token: Token): Token
  state: State
}

export function useStore(): Store | undefined {
  return inject(storeSymbol);
}
