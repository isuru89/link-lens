const HEADING_LABELS = {
  H1: "H1",
  H2: "H2",
  H3: "H3",
  H4: "H4",
  H5: "H5",
  H6: "H6",
};

const Chip = ({ label, value, bgColor, valueBgColor, color }) => {
  return (
    <div
      className="chip"
      style={{ backgroundColor: bgColor, borderColor: color, color }}
    >
      <div className="chip-label">{label}</div>
      <div className="chip-value" style={{ backgroundColor: valueBgColor }}>
        {value}
      </div>
    </div>
  );
};

const LinkStatsSection = ({
  InternalLinkCount = 0,
  ExternalLinkCount = 0,
  InvalidLinkCount = 0,
  InvalidLinks = [],
}) => {
  return (
    <>
      <DataRow
        label={"Links:"}
        value={
          <div style={{ display: "flex", gap: "4px" }}>
            <Chip label={"Internal"} value={InternalLinkCount} />
            <Chip label={"External"} value={ExternalLinkCount} />
            <Chip
              label={"Invalid"}
              value={InvalidLinkCount}
              bgColor={"#ff000011"}
              valueBgColor={"#ff000022"}
              color={"#ff0000"}
            />
          </div>
        }
      />
      {InvalidLinks && InvalidLinks.length > 0 && (
        <DataRow
          id="invalid-links"
          label={"Invalid Links:"}
          value={
            <div>
              {InvalidLinks.map((l) => (
                <div>• {l}</div>
              ))}
            </div>
          }
        />
      )}
    </>
  );
};

const HeadingsSection = ({ HeadingsCount }) => {
  return (
    <DataRow
      label={"Headings:"}
      value={
        <div style={{ display: "flex", gap: "4px" }}>
          {Object.keys(HeadingsCount).map((key) => {
            return (
              <Chip
                label={`${HEADING_LABELS[key] || key}`}
                value={HeadingsCount[key] || 0}
              />
            );
          })}
        </div>
      }
    />
  );
};

const DataRow = ({ id, label, value }) => {
  return (
    <div id={id} className="datarow">
      <div className="datarow-label">{label}</div>
      <div className="datarow-value">{value}</div>
    </div>
  );
};

const formatElapsed = (elapsed) => {
  return Math.fround((elapsed / 1000.0) * 1000) / 1000;
};

export const LinkInfo = ({ data, elapsed = -1 }) => {
  const {
    HtmlVersion,
    Title,
    HeadingsCount = {},
    LinkStats = {},
    PageType,
  } = data;

  return (
    <div className="result-panel">
      {elapsed > 0 && (
        <div
          style={{ textAlign: "right", fontSize: "12px", fontStyle: "italic" }}
        >
          Elapsed: {formatElapsed(elapsed)} seconds
        </div>
      )}
      <DataRow label={"HTML Version:"} value={HtmlVersion} />
      <DataRow label={"Title:"} value={Title} />
      <DataRow label={"Page Type:"} value={PageType} />
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
