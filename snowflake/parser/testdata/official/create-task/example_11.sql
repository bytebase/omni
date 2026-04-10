CREATE TASK finalize_task
  WAREHOUSE = my_warehouse
  FINALIZE = my_root_task
  AS
    CALL SYSTEM$SEND_EMAIL(
      'my_email_int',
      'first.last@example.com, first2.last2@example.com',
      'Email Alert: Task A has finished.',
      'Task A has successfully finished.\nStart Time: 10:10:32\nEnd Time: 12:15:45\nTotal Records Processed: 115678'
    );
