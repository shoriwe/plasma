#include "vm/virtual_machine.h"


plasma::vm::constructor_callback plasma::vm::virtual_machine::runtime_error_initialize(bool isBuiltIn) {
    return [this, isBuiltIn](context *c, struct value *object) -> value * {
        object->set_on_demand_symbol(
                Initialize,
                [this, c, object, isBuiltIn]() -> value * {
                    return this->new_function(
                            c, isBuiltIn,
                            object,
                            new_builtin_callable(
                                    1,
                                    [this, c](value *self, const std::vector<value *> &arguments,
                                              bool *success) -> value * {
                                        value *message = arguments[0];
                                        if (message->typeId != String) {
                                            (*success) = false;
                                            return this->new_invalid_type_error(
                                                    c,
                                                    message->get_type(c, this),
                                                    std::vector<std::string>{StringName}
                                            );
                                        }
                                        self->string = message->string;
                                        (*success) = true;
                                        return this->get_none(c);
                                    }
                            )
                    );
                }
        );
        object->set_on_demand_symbol(
                ToString,
                [this, c, object, isBuiltIn]() -> value * {
                    return this->new_function(
                            c, isBuiltIn,
                            object,
                            new_builtin_callable(
                                    0,
                                    [this, c](value *self, const std::vector<value *> &arguments,
                                              bool *success) -> value * {

                                        (*success) = true;
                                        return this->new_string(
                                                c, false,
                                                self->get_type(c, this)->name + ": " + self->string
                                        );
                                    }
                            )
                    );
                }
        );
        return nullptr;
    };
}