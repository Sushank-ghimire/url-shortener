"use client";

import { useForm } from "@tanstack/react-form";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { signupFormSchema } from "@/schemas/auth";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { useRegister } from "@/lib/query/actions/auth/hooks/use-register";
import { toast } from "sonner";

export default function SignupFrom() {
  const registerMutation = useRegister();

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
      name: "",
    },
    validators: {
      onSubmit: signupFormSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        await registerMutation.mutateAsync(value, {
          onSuccess: (data) => {
            if (data.success) {
              toast.success("account created successfully!", {
                description: "verify your email to get started",
              });
              form.reset();
            }
          },
        });
      } catch (error) {
        toast.error(
          error instanceof Error ? error.message : "Invalid email or password.",
        );
      }
    },
  });

  const { state } = form;
  return (
    <Card className="sm:w-full p-4 md:w-xl">
      <CardHeader className="w-full items-center justify-center font-bold lg:text-xl text-lg">
        Let's get started
      </CardHeader>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
      >
        <CardContent>
          <FieldGroup>
            <form.Field
              name="name"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;

                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name}>Full Name</FieldLabel>

                    <Input
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      placeholder="Sushank Ghimire"
                    />

                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />

            <form.Field
              name="email"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;

                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name}>Email</FieldLabel>

                    <Input
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      placeholder="mail@gmail.com"
                      autoComplete="email"
                      type="email"
                    />

                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />

            <form.Field
              name="password"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;

                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name}>Password</FieldLabel>

                    <Input
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      placeholder="********"
                      type="password"
                    />

                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
          </FieldGroup>
        </CardContent>

        <CardFooter className="mt-4 w-full justify-center items-center flex flex-col gap-4">
          <Button
            disabled={state.isSubmitting || state.errors.length > 0}
            className={"cursor-pointer w-full"}
            type="submit"
          >
            {form.state.isSubmitting ? "Creating account..." : "Register"}
          </Button>
          <div className="flex gap-2">
            <p>Already have an account login here ?</p>
            <Link className="font-bold" href={"/auth/login"}>
              Login
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
}
