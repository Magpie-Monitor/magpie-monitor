import './Table.scss';

interface TableRow {
  [key: string]: string | number | boolean;
}

export interface TableProps<Row extends TableRow> {
  columns: TableColumn<Row>[];
  rows: Row[];

  /**
   * Max height of the table defined in the css units (px, rem, ect.)
   * Used to accomodate sticky table headers.
   * By default table doesn't have a max height.
   */
  maxHeight?: string;

  /**
   * AlignLeft forces all columns to take as little space as possible
   * making all columns tightly packed on the left side of the table.
   * Useful if we need a column that is always the size of its contents.
   * Default value is false
   */
  alignLeft?: boolean;
}

export interface TableColumn<Row extends TableRow> {
  header: string;
  columnKey: keyof Row;
  customComponent?: (row: Row) => React.ReactNode;
}
export interface TableRowProps<Row extends TableRow> {
  row: Row;
  columns: TableColumn<Row>[];
  alignLeft: boolean;
}

export const getCellRowSpan = <Row extends TableRow>(
  columns: TableColumn<Row>[],
  alignLeft: boolean,
): { gridColumn: string } => {
  if (columns.length == 0 || alignLeft) {
    return { gridColumn: '' };
  }

  return { gridColumn: `span ${Math.floor(100 / columns.length)}` };
};

function getTableStyle(maxHeight: string) {
  return { maxHeight };
}

function TableBodyRow<Row extends TableRow>({
  row,
  columns,
  alignLeft,
}: TableRowProps<Row>): React.ReactNode {
  return (
    <div className="table__row">
      {columns.map((column: TableColumn<Row>) => (
        <div key="index" className={'table__cell'} style={getCellRowSpan(columns, alignLeft)}>
          {column.customComponent ? column.customComponent(row) : row[column.columnKey]}
        </div>
      ))}
    </div>
  );
}

function Table<T extends TableRow>({
  columns,
  rows,
  maxHeight,
  alignLeft = false,
}: TableProps<T>): React.ReactNode {
  const tableStyle = maxHeight ? getTableStyle(maxHeight) : {};

  return (
    <div className="table">
      <div className="table__headers">
        {columns.map((column: TableColumn<T>, index: number) => (
          <div key={index} className="table__header" style={getCellRowSpan(columns, alignLeft)}>
            {column.header}
          </div>
        ))}
      </div>
      <div className="table__body" style={tableStyle}>
        {rows.map((row: T, index: number) => (
          <TableBodyRow key={index} row={row} columns={columns} alignLeft={alignLeft} />
        ))}
      </div>
    </div>
  );
}

export default Table;
