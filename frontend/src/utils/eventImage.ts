import placeholder from '@/assets/event-placeholder.svg'

/** Returns the event thumbnail URL, falling back to the branded placeholder. */
export function eventImage(url: string | undefined | null): string {
  return url || placeholder
}

/** @error handler — swaps a broken <img> src to the placeholder. */
export function onImageError(e: Event) {
  const img = e.target as HTMLImageElement
  if (img.src !== placeholder) {
    img.src = placeholder
  }
}

export { placeholder as EVENT_PLACEHOLDER }
