<script setup lang="ts">
import {
  useSlideContext,
  onSlideEnter,
} from '@slidev/client'
import { motion, AnimatePresence } from 'motion-v'
import { ref, watch } from 'vue'
import Fiber from "./Fiber.vue"
import { Icon } from "@iconify/vue";

let ctx = useSlideContext()

type Fiber = {
  id: number,
  name: string,
  blocked: boolean,
}

const fibersFromToo = (start: number, end: number): Fiber[] => {
  return Array.from({ length: end - start + 1 }).map((_, i) => {
    return fiber(start + i)
  })
}

const fiber = (id: number, blocked: boolean = false): Fiber => {
  return {
    id,
    name: `Fiber ${id + 1}`,
    blocked,
  }
}

type Core = {
  id: number,
  name: string,
  fiber: Fiber | null,
}

const coreFromToo = (start: number, end: number): Core[] => {
  return Array.from({ length: end - start + 1 }).map((_, i) => {
    return core(start + i)
  })
}

const core = (id: number, fiber: Fiber | null = null): Core => {
  return {
    id,
    name: `Core ${id + 1}`,
    fiber,
  }
}


const fibersInQueue = ref<Fiber[]>()

const fibersCompleted = ref<Fiber[]>()

const cores = ref<Core[]>()

const initialize = () => {
  fibersInQueue.value = fibersFromToo(0, 7)
  cores.value = coreFromToo(0, 3)
  fibersCompleted.value = []
}

onSlideEnter(() => {
  if (ctx.$clicks.value === 0) {
    initialize()
  }
})

watch(ctx.$clicks, (clicks, oldClicks) => {
  switch (clicks) {
    case 1:
      fibersInQueue.value = [
        fiber(1),
        fiber(2),
        fiber(3),
        fiber(4),
        fiber(5),
        fiber(6),
        fiber(7),
      ]
      cores.value = [
        core(0, fiber(0)),
        core(1),
        core(2),
        core(3),
      ]
      fibersCompleted.value = []
      break
    case 2:
      fibersInQueue.value = [
        fiber(2),
        fiber(3),
        fiber(4),
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1)),
        core(2),
        core(3),
      ]
      fibersCompleted.value = []
      break
    case 3:
      fibersInQueue.value = [
        fiber(3),
        fiber(4),
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1)),
        core(2, fiber(2)),
        core(3),
      ]
      fibersCompleted.value = []
      break
    case 4:
      fibersInQueue.value = [
        fiber(4),
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1)),
        core(2, fiber(2)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = []
      break
    case 5:
      fibersInQueue.value = [
        fiber(4),
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1)),
        core(2),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 6:
      fibersInQueue.value = [
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1)),
        core(2, fiber(4)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 7:
      fibersInQueue.value = [
        fiber(5),
        fiber(6),
        fiber(7),

      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(1, true)),
        core(2, fiber(4)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 8:
      fibersInQueue.value = [
        fiber(5),
        fiber(6),
        fiber(7),
        fiber(1, true)
      ]
      cores.value = [
        core(0, fiber(0)),
        core(1),
        core(2, fiber(4)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 9:
      fibersInQueue.value = [
        fiber(6),
        fiber(7),
        fiber(1, true)
      ]
      cores.value = [
        core(0, fiber(0)),
        core(1, fiber(5)),
        core(2, fiber(4)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 10:
      fibersInQueue.value = [
        fiber(6),
        fiber(7),
        fiber(1, true)
      ]
      cores.value = [
        core(0, fiber(0, true)),
        core(1, fiber(5)),
        core(2, fiber(4)),
        core(3, fiber(3, true)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 11:
      fibersInQueue.value = [
        fiber(6),
        fiber(7),
        fiber(1, true),
        fiber(0, true),
        fiber(3, true)
      ]
      cores.value = [
        core(0),
        core(1, fiber(5)),
        core(2, fiber(4)),
        core(3),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 12:
      fibersInQueue.value = [
        fiber(1, true),
        fiber(0, true),
        fiber(3, true)
      ]
      cores.value = [
        core(0, fiber(6)),
        core(1, fiber(5)),
        core(2, fiber(4)),
        core(3, fiber(7)),
      ]
      fibersCompleted.value = [fiber(2)]
      break
    case 13:
      fibersInQueue.value = [
        fiber(1, true),
        fiber(0, true),
        fiber(3, true)
      ]
      cores.value = [
        core(0, fiber(6)),
        core(1),
        core(2),
        core(3, fiber(7)),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4)]
      break
    case 14:
      fibersInQueue.value = [
        fiber(3, true)
      ]
      cores.value = [
        core(0, fiber(6)),
        core(1, fiber(1)),
        core(2, fiber(0)),
        core(3, fiber(7)),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4)]
      break
    case 15:
      fibersInQueue.value = [
        fiber(3, true)
      ]
      cores.value = [
        core(0, fiber(6)),
        core(1, fiber(1)),
        core(2, fiber(0)),
        core(3),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4), fiber(7)]
      break
    case 16:
      fibersInQueue.value = [

      ]
      cores.value = [
        core(0, fiber(6)),
        core(1, fiber(1)),
        core(2, fiber(0)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4), fiber(7)]
      break
    case 17:
      fibersInQueue.value = [

      ]
      cores.value = [
        core(0),
        core(1),
        core(2, fiber(0)),
        core(3, fiber(3)),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4), fiber(7), fiber(6), fiber(1)]
      break
    case 18:
      fibersInQueue.value = [

      ]
      cores.value = [
        core(0),
        core(1),
        core(2),
        core(3),
      ]
      fibersCompleted.value = [fiber(2), fiber(5), fiber(4), fiber(7), fiber(6), fiber(1), fiber(0), fiber(3)]
      break

    default:
      initialize()
      break
  }
})

</script>

<template>
  <div class="flex flex-row gap-2 position-fixed top-0">
    <div class="flex flex-row gap-2 items-center">
      <div class="h-[8px] w-[8px] bg-accent" />
      <p class="text-size-xs!">Queued</p>
    </div>
    <div class="flex flex-row gap-2 items-center">
      <div class="h-[8px] w-[8px] bg-blue" />
      <p class="text-size-xs!">Running</p>
    </div>
    <div class="flex flex-row gap-2 items-center">
      <div class="h-[8px] w-[8px] bg-green" />
      <p class="text-size-xs!">Finished</p>
    </div>
    <div class="flex flex-row gap-2 items-center">
      <div class="h-[8px] w-[8px] bg-yellow" />
      <p class="text-size-xs!">Waiting on I/O</p>
    </div>
  </div>

  <div class="flex flex-row w-full h-full justify-between">

    <div class="b-4 b-accent flex flex-col p-2 h-full w-[200px]">
      <h3 class="text-white">Queue</h3>
      <div class="flex flex-col h-full">
        <motion.div layout v-for="fiber in fibersInQueue" :key="fiber.id">
          <AnimatePresence>
            <Fiber status="queued" :fiber="fiber" />
          </AnimatePresence>
        </motion.div>
      </div>
    </div>

    <div class="b-4 b-accent flex flex-col p-2 h-full w-[200px]">
      <h3 class="text-white">Running</h3>
      <div layout class="flex flex-col h-full justify-between">
        <motion.div layout v-for="core in cores" :key="core.id">
          <div class="b-b-coolgray b-b-2 p-y-2">
            <p class="m-0! flex flex-row gap-2 items-center">
              {{ core.name }}
              <Icon icon="solar:cpu-bold" width="16" height="16" />
            </p>
            <Fiber v-if="core.fiber" status="running" :fiber="core.fiber" />
            <div v-else class="h-[50px]" />
          </div>
        </motion.div>
      </div>
    </div>

    <div class="b-4 b-accent flex flex-col p-2 h-full w-[200px]">
      <h3 class="text-white">Finished</h3>
      <div class="flex flex-col h-full">
        <motion.div layout v-for="fiber in fibersCompleted" :key="fiber.id">
          <AnimatePresence>
            <Fiber status="finished" :fiber="fiber" />
          </AnimatePresence>
        </motion.div>
      </div>
    </div>
  </div>
</template>
