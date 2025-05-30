<script setup>
  const items = ref([
    {
      title: 'Константы',
      isOpen: true,
      children: [
        {
          title: 'OP_0',
          value: 'OP_0',
        },
        {
          title: 'OP_1',
          value: 'OP_1',
        },
        {
          title: 'OP_2',
          value: 'OP_2',
        },
      ],
    },
    {
      title: 'Работа со стеком',
      isOpen: true,
      children: [

        {
          title: 'OP_PUSHDATA',
          value: 'OP_PUSHDATA',
        },
        {
          title: 'OP_DUP',
          value: 'OP_DUP',
        },
        {
          title: 'OP_DROP',
          value: 'OP_DROP',
        },
      ],
    },
  ])

  const open = ref([])
</script>

<template>
  <v-card width="250">
    <v-list v-model:opened="open" density="compact">
      <v-list-group
        v-for="(group, index) in items"
        :key="index"
        :title="group.title"
        :value="group.title"
      >
        <template #activator="{ props }">

          <v-list-item
            v-bind="props"
            :title="group.title"
          />
        </template>
        <v-list-item
          v-for="(item, i) in group?.children"
          :key="i"
          :title="item.title"
          :value="item.value"
        />
      </v-list-group>
    </v-list>
  </v-card>
</template>
