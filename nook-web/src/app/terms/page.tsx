export const metadata = {
  title: "Terms of Service | Nook",
};

export default function Terms() {
  return (
    <div>
      <div className="prose prose-invert mx-auto my-20 px-3 md:px-5">
        <h1>Terms of Service</h1>
        <p>Last updated: September 23, 2025</p>
        <div>
          <p>
            By inviting Nook to your Discord server or logging into our website
            (https://nooks.chat) you agree that you have read, understood, and
            accepted these terms. You are also responsible for informing the
            members in your Discord server about these terms.
          </p>
          <p>
            If you do not agree with any of these terms, you are prohibited from
            using or adding any version of Nook to your server.
          </p>

          <h2>Disclaimer</h2>
          <p>
            You are strictly prohibited to use Nook against the ToS of Discord
            or for illegal purposes. We are doing our best to prevent these
            activities, while trying to provide the best user experience as
            possible. If you find people or communities using Nook against the
            ToS of discord or even for illegal activities, please send us an
            email to{" "}
            <span className="text-blue-400">contact (at) nooks (dot) chat</span>
            .
          </p>

          <h2>Proprietary Rights</h2>
          <p>
            We (Nook or more specifically Merlin Fuchs) own and retain all
            rights for public available data (including but not limited to
            templates). We grant you the permission to use this available data
            for your own needs, but strictly disallow any commercial use. You
            therefore shall not sell, license or otherwise commercialize the
            data except if the permission was expressly granted to you.
          </p>

          <h2>Availability</h2>
          <ul>
            <li>
              Nook is provided as-is. There are no guarantees that it will be
              available in the future, and its purpose or availability may be
              changed at any time.
            </li>
            <li>
              User related data including backups may be deleted at any time.
            </li>
            <li>
              User related data including backups is non-transferable between
              discord accounts.
            </li>
            <li>
              Any premium features are not guaranteed. They may change or be
              revoked at any time.
            </li>
            <li>
              Access to all or specific features of Nook may be revoked, for all
              or a specific user, at any time.
            </li>
          </ul>

          <h2>Refund Policy</h2>
          <p>
            Refunds for payments may be issued when the user requests it and one
            of the following requirements is met:
          </p>
          <ul>
            <li>
              it&apos;s the first payment and the user requests a refund on the
              same day
            </li>
            <li>
              it&apos;s a monthly recurring payment and the user requests a
              refund in a three day period
            </li>
            <li>
              it&apos;s an annual recurring payment and the user requests a
              refund in a seven day period
            </li>
            <li>
              A major premium feature was removed and the user bought the tier
              specifically because of that feature
            </li>
          </ul>
          <p>
            We still require a comprehensible explanation of the situation and
            never guarantee a refund. Partial refunds might be issued for annual
            payments under special circumstances. Refunds are never issued if
            you have broken our Terms of Service and therefore lost access to
            certain or all features.
          </p>
        </div>
      </div>
    </div>
  );
}
