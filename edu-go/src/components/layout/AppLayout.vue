<script setup lang="ts">
import AppAside from './AppAside.vue';
</script>


<template>
    <template>
        <div class="common-layout">
            <el-container>
                <AppAside />
                <el-container>
                    <el-header>Header</el-header>
                    <el-main>Main</el-main>
                </el-container>
            </el-container>
        </div>
    </template>
</template>


<style scoped></style>