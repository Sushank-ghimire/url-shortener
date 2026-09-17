import RedirectPage from "@/components/redirect";

export default async function RedirectUser({
  params,
}: {
  params: Promise<{ shortcode: string }>;
}) {
  const { shortcode } = await params;
  return <RedirectPage shortcode={shortcode} />;
}
