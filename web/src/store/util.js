import { MESSAGE_ERROR, } from "../constants";

export const call = ({ commit, dispatch, promise, loadingMutation }) => new Promise((resolve, reject) => {
  promise
    .then(res => {
      if (!res.ok) {
        dispatch("addMessage", { type: MESSAGE_ERROR, text: res.error });
        reject(res)
      } else {
        resolve(res);
      }
    })
    .catch(error => {
      dispatch("addMessage", { type: MESSAGE_ERROR, text: error.message });
      reject({ ok: false, code: 0, error: error.message });
    })

  if (loadingMutation) {
    commit(loadingMutation, true);
    promise.finally(() => {
      commit(loadingMutation, false);
    })
  }
})

export const radioCall = (context) => {
  return new Promise((resolve, reject) => {
    call(context)
      .then(res => resolve(res))
      .catch((res) => {
        if (res.code == 404) {
          context.dispatch("listRadios")
        }
        reject(res)
      })
  })
}