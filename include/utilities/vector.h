#ifndef VECTOR_H
#define VECTOR_H

#include <stdlib.h>
#include <stdio.h>

#define VECTOR_INIT_CAPACITY 8
#define VECTOR_CAPACITY_SCALER 2

typedef struct {
    void** data;
    size_t size;
    size_t capacity;
    void (*elem_destructor)(void*);
} Vector;

static inline void vector_init_with_capacity(Vector* vector, void (*elem_destructor)(void*), size_t capacity) {
    vector->capacity = capacity;
    vector->size = 0;
    vector->data = (void**)malloc(vector->capacity * sizeof(void*));

    if (vector->data == NULL) {
        printf("Vector memory allocation failed!\n");
        exit(1);
    }

    vector->elem_destructor = elem_destructor;
}

static inline Vector* vector_create_with_capacity(size_t capacity, void (*elem_destructor)(void*)) {
    Vector* vector = malloc(sizeof(Vector));
    if (vector == NULL) {
        printf("Vector memory allocation failed!\n");
        exit(1);
    }
    vector_init_with_capacity(vector, elem_destructor, capacity);

    return vector;
}

static inline Vector* vector_create(void (*elem_destructor)(void*)) {
    return vector_create_with_capacity(VECTOR_INIT_CAPACITY, elem_destructor);
}

static inline void vector_init(Vector* vector, void (*elem_destructor)(void*)) {
    vector_init_with_capacity(vector, elem_destructor, VECTOR_INIT_CAPACITY);
}

static inline void vector_resize(Vector* vector, size_t new_capacity) {
    void** temp = (void**)realloc(vector->data, new_capacity * sizeof(void*));
    if (temp == NULL) {
        printf("Vector resize failed!\n");
        exit(1);
    }

    vector->data = temp;
    vector->capacity = new_capacity;
}

static inline void vector_push(Vector* vector, void* elem) {
    if (vector->size >= vector->capacity) {
        vector_resize(vector, vector->capacity * VECTOR_CAPACITY_SCALER);
    }

    vector->data[vector->size++] = elem;
}

static inline void* vector_get(Vector* vector, size_t index) {
    if (index >= vector->size) {
        return NULL;
    }

    return vector->data[index];
}

static inline void vector_pop(Vector* vector) {
    if (vector->size > 0) {
        if (vector->elem_destructor) {
            vector->elem_destructor(vector->data[vector->size - 1]);
        }
        --vector->size;
    }
}

static inline void vector_free(Vector* vector) {
    if (vector) {
        if (vector->elem_destructor) {
            for (size_t i = 0; i < vector->size; ++i) {
                vector->elem_destructor(vector->data[i]);
            }
        }
        free(vector->data);
        vector->data = NULL;
        vector->size = 0;
        vector->capacity = 0;
        free(vector);
    }
}

static inline void* vector_back(Vector* vector) {
    if (vector->size > 0) {
        return vector->data[vector->size - 1];
    }

    return NULL;
}

#endif // VECTOR_H
