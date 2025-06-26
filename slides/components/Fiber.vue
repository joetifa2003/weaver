<script lang="ts" setup>
import { motion } from 'motion-v'
import { ref } from 'vue'
import { Icon } from "@iconify/vue";

const props = defineProps<{
  fiber: {
    id: number,
    name: string,
    blocked: boolean,
  },
  status: "queued" | "running" | "blocked" | "finished" | "blocked",
}>()

const statusClasses = {
  queued: "bg-accent",
  running: "bg-blue",
  finished: "bg-green",
  blocked: "bg-yellow",
}

</script>

<template>
  <motion.div :layoutId="`${props.fiber.id}`"
    :class="`flex p-4 items-center justify-center h-[50px] gap-2 ${props.fiber.blocked ? statusClasses['blocked'] : statusClasses[props.status]}`"
    :transition="{
      duration: 0.4,
      ease: 'easeInOut',
    }">
    {{ props.fiber.name }}
    <Icon v-if="props.status === 'running'" icon="eos-icons:arrow-rotate" width="16" height="16" />
    <Icon v-if="props.status === 'finished'" icon="fa-solid:check" width="16" height="16" />
    <Icon v-if="props.status === 'blocked'" icon="svg-spinners:blocks-shuffle-3" width="16" height="16" />
    <Icon v-if="props.status === 'queued'" icon="eos-icons:hourglass" width="16" height="16" />
    <Icon v-if="props.fiber.blocked" icon="svg-spinners:blocks-shuffle-3" width="16" height="16" />
  </motion.div>
</template>
