async function fetch_const() {
  const response = await fetch('http://localhost:3000/v1/const', {
    headers: {
      Accept: "application/json",
    },
    method: 'GET'
  });

  if (!response.ok) {
    throw Error(`${response.status} : unable to fetch const values`);
  }

  const const_json = await response.json();


  console.log('const data fetched.');
  return const_json;
}

async function fetch_serialized_schedule(selected_semester) {
  const response = await fetch(`http://localhost:3000/v1/university_schedule?semester=${selected_semester}`, {
    headers: {
      Accept: "text/plain",
    },
    method: 'GET'
  });

  if (!response.ok) {
    const err_msg = await response.text();
    throw Error(`${response.status} : ${err_msg}`);
  }

  const serialized_schedule = await response.arrayBuffer();

  return [new Uint8Array(serialized_schedule), response.ok];
}

async function fetch_serialized_class_schedule(department_id, selected_semester, schedule_idx) {
  const response = await fetch(
    `http://localhost:3000/v1/class_schedule?department_id=${department_id}&semester=${selected_semester}&schedule_idx=${schedule_idx}`, {
    headers: {
      Accept: "text/plain",
    },
    method: 'GET'
  }
  );

  if (!response.ok) {
    const err_msg = await response.text();
    throw Error(`${response.status} : ${err_msg}`);
  }

  const serialized_schedule = await response.arrayBuffer();

  return [new Uint8Array(serialized_schedule), response.ok];
}

async function send_serialized_schedule(selected_semester, serialized_schedule) {
  const response = await fetch(`http://localhost:3000/v1/university_schedule?semester=${selected_semester}`, {
    method: 'POST',
    headers: {
      Accept: "text/plain",
      'Content-Type': 'application/octet-stream'
    },
    body: serialized_schedule
  });


  if (!response.ok) {
    const err_msg = await response.text();
    throw new Error(`${response.status} : ${err_msg}`);
  }

  console.log('serialized university schedule sent to backend.');
}

async function deserialize_schedule(serialized_data) {
  let constants = await fetch_const()

  const time_slot_bytes = constants.time_slot_bytes;
  const weekly_school_days = constants.weekly_school_days;

  const daily_time_slots = constants.daily_time_slots;
  const weekly_time_slots = constants.weekly_time_slots;

  const university_schedules = [];
  const number_of_sections = serialized_data.length / (weekly_time_slots * time_slot_bytes);

  let offset = 0;

  for (let section_idx = 0; section_idx < number_of_sections; section_idx++) {
    const week_time_table = [];
    for (let day = 0; day < weekly_school_days; day++) {
      const day_time_table = [];
      for (let time_slot = 0; time_slot < daily_time_slots; time_slot++) {
        const subjectID = (serialized_data[offset] | (serialized_data[offset + 1] << 8));
        const instructorID = (serialized_data[offset + 2] | (serialized_data[offset + 3] << 8));
        const roomID = (serialized_data[offset + 4] | (serialized_data[offset + 5] << 8));

        day_time_table.push({ subjectID, instructorID, roomID });
        offset += time_slot_bytes;
      }
      week_time_table.push(day_time_table);
    }
    university_schedules.push(week_time_table);
  }

  console.log('university schedule deserialized.');
  return university_schedules;
}

async function serialize_schedule(university_schedules) {
  let constants = await fetch_const()

  const time_slot_bytes = constants.time_slot_bytes;
  const weekly_time_slots = constants.weekly_time_slots;

  const serialized_data = new Uint8Array(university_schedules.length * weekly_time_slots * time_slot_bytes);

  let offset = 0;

  for (const section_week_schedule of university_schedules) {
    for (const day of section_week_schedule) {
      for (const time_slot of day) {
        serialized_data[offset] = time_slot.subjectID & 0xff;
        serialized_data[offset + 1] = (time_slot.subjectID >> 8) & 0xff;
        serialized_data[offset + 2] = time_slot.instructorID & 0xff;
        serialized_data[offset + 3] = (time_slot.instructorID >> 8) & 0xff;
        serialized_data[offset + 4] = time_slot.roomID & 0xff;
        serialized_data[offset + 5] = (time_slot.roomID >> 8) & 0xff;

        offset += time_slot_bytes;
      }
    }
  }

  console.log('university schedule serialized.');
  return serialized_data;
}

async function test() {
  try {
    let [raw_data, _] = await fetch_serialized_schedule(0);
    let deserialized = await deserialize_schedule(raw_data);
    let serialized = await serialize_schedule(deserialized);
    await send_serialized_schedule(0, serialized)
  
    await fetch(`http://localhost:3000/die`, {
      headers: {
        Accept: "text/plain",
      },
      method: 'GET'
    });
  } catch (err) {
    if (err.cause?.code === 'UND_ERR_SOCKET') {
      console.log('err.cause.code =', err.cause.code)
      process.exit(0)
    } else {
      console.log('err =', err)

      await fetch(`http://localhost:3000/die`, {
        headers: {
          Accept: "text/plain",
        },
        method: 'GET'
      });

      process.exit(1)
    }
  }
}

test()