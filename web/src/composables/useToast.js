import { ref } from 'vue'

const notifications = ref([])
let idCounter = 0

export function useToast() {
  const show = (message, type = 'info', duration = 3000) => {
    const id = ++idCounter
    notifications.value.push({
      id,
      message,
      type,
      visible: true
    })

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }

    return id
  }

  const remove = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value[index].visible = false
      setTimeout(() => {
        notifications.value.splice(index, 1)
      }, 300)
    }
  }

  const success = (message, duration = 3000) => show(message, 'success', duration)
  const error = (message, duration = 4000) => show(message, 'error', duration)
  const warning = (message, duration = 3000) => show(message, 'warning', duration)
  const info = (message, duration = 3000) => show(message, 'info', duration)

  return {
    notifications,
    show,
    remove,
    success,
    error,
    warning,
    info
  }
}
