


export function copyObject(item: any) {
  return JSON.parse(JSON.stringify(item))
}

export function findArrayIndex(array: any, col: string, wert: any, col2?: string, wert2?: any) {

    let lReturn = -1;
  
    for (let i = 0; i < array.length; i++) {
  
        if (col2 && col2 != undefined && wert2 && wert2 != undefined) {
            if (array[i][col] && array[i][col] === wert && array[i][col2] && array[i][col2] === wert2) {
                lReturn = i;
                break;
            }
        } else {
            if (array[i][col] && array[i][col] === wert) {
                lReturn = i;
                break;
            }
        }
    }
    return lReturn;
  }

  export const debounce = (func: any, delay = 600, immediate = false) => {
    let timeout: any
    return function () {
    //@ts-ignore
      const context = this
      const args = arguments
      const later = function () {
        timeout = null
        if (!immediate) func.apply(context, args)
      }
      const callNow = immediate && !timeout
      clearTimeout(timeout)
      timeout = setTimeout(later, delay)
      if (callNow) func.apply(context, args)
    }
  }