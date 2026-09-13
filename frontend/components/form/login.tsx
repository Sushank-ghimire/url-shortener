"use client";

import { useForm } from "@tanstack/react-form";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { loginFormSchema } from "@/schemas/auth";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { useLogin } from "@/lib/query/actions/auth/hooks/use-login";
import { toast } from "sonner";
import secureLocalStorage from "react-secure-storage";
import { useRouter } from "next/navigation";

export default function LoginFrom() {
  const loginMutation = useLogin();

  const router = useRouter();

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
    },
    validators: {
      onSubmit: loginFormSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        await loginMutation.mutateAsync(value, {
          onSuccess: ({ data, success }) => {
            if (success) {
              secureLocalStorage.setItem("refresh_token", data.refresh_token);
              secureLocalStorage.setItem("access_token", data.access_token);
              toast.success("Login successful");
              router.replace("/user/dashboard");
            }
          },
        });
        form.reset();
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
        Login to continue
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
            {state.isSubmitting ? "Logging in..." : "Login"}
          </Button>
          <div className="flex gap-2">
            <p>Don't have an account register here ?</p>
            <Link className="font-bold" href={"/auth/register"}>
              Register
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
}
