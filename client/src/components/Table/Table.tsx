import './Table.scss';

interface TableRow {
  [key: string]: string | number | boolean;
}

export interface TableProps<Row extends TableRow> {
  columns: TableColumn<Row>[];
  rows: Row[];
}

export interface TableColumn<Row extends TableRow> {
  header: string;
  columnKey: keyof Row;
  customComponent?: (row: Row) => React.ReactNode;
}
export interface TableRowProps<Row extends TableRow> {
  row: Row;
  columns: TableColumn<Row>[];
}

export const getCellRowSpan = <Row extends TableRow>(
  columns: TableColumn<Row>[],
): { gridColumn: string } => {
  return { gridColumn: `span ${columns.length}` };
};

function TableBodyRow<Row extends TableRow>({ row, columns }: TableRowProps<Row>): React.ReactNode {
  return (
    <div className="table__row">
      {columns.map((column: TableColumn<Row>) => (
        <div key="index" className={'table__cell'}>
          {column.customComponent ? column.customComponent(row) : row[column.columnKey]}
        </div>
      ))}
    </div>
  );
}

function Table<T extends TableRow>({ columns, rows }: TableProps<T>): React.ReactNode {
  return (
    <div className="table">
      <div className="table__headers">
        {columns.map((column: TableColumn<T>, index: number) => (
          <div key={index} className="table__header">
            {column.header}
          </div>
        ))}
      </div>
      <div className="table__body">
        {rows.map((row: T, index: number) => (
          <TableBodyRow key={index} row={row} columns={columns} />
        ))}
      </div>
    </div>
  );
}

export default Table;
