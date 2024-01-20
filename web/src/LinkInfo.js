const HEADING_LABELS = {
  H1: "Heading-1",
  H2: "Heading-2",
  H3: "Heading-3",
  H4: "Heading-4",
  H5: "Heading-5",
  H6: "Heading-6",
};

const LinkStatsSection = ({
  InternalLinkCount = 0,
  ExternalLinkCount = 0,
  InvalidLinkCount = 0,
  InvalidLinks = [],
}) => {
  return (
    <div>
      <span>Links: </span>
      <div>
        <div>
          <span>Internal Links: </span>
          <span>{InternalLinkCount}</span>
        </div>
        <div>
          <span>External Links: </span>
          <span>{ExternalLinkCount}</span>
        </div>
        <div>
          <span>Inaccessible Links: </span>
          <span>{InvalidLinkCount}</span>
        </div>
        {InvalidLinks && InvalidLinks.length > 0 && (
          <ol>
            {InvalidLinks.map((lnk) => (
              <li>{lnk}</li>
            ))}
          </ol>
        )}
      </div>
    </div>
  );
};

const HeadingsSection = ({ HeadingsCount }) => {
  return (
    <div>
      <span>Headings Count: </span>
      <div>
        {Object.keys(HeadingsCount).map((key) => {
          return (
            <div>
              <span>{HEADING_LABELS[key] || key}: </span>
              <span>{HeadingsCount[key] || 0}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export const LinkInfo = ({ data }) => {
  const {
    HtmlVersion,
    Title,
    HeadingsCount = {},
    LinkStats = {},
    PageType,
  } = data;

  return (
    <div>
      <div>
        <span>Html Version: </span>
        <span>{HtmlVersion}</span>
      </div>
      <div>
        <span>Title: </span>
        <span>{Title}</span>
      </div>
      <div>
        <span>Page Type: </span>
        <span>{PageType}</span>
      </div>
      <HeadingsSection HeadingsCount={HeadingsCount} />
      <LinkStatsSection
        InternalLinkCount={LinkStats["InternalLinkCount"]}
        ExternalLinkCount={LinkStats["ExternalLinkCount"]}
        InvalidLinkCount={LinkStats["InvalidLinkCount"]}
        InvalidLinks={LinkStats["InvalidLinks"]}
      />
    </div>
  );
};
