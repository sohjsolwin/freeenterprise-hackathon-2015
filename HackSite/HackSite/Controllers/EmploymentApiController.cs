using HackSite.Dto;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Web.Http;

namespace HackSite.Controllers
{
    [RoutePrefix("api/employmentdata")]
    public class EmploymentApiController : ApiController
    {

        private static IEnumerable<UnemploymentData> cachedUnEmpData;


        [Route("getData")]
        public IHttpActionResult GetGraphData([FromUri] DateTime date)
        {
            var endDate = new DateTime(DateTime.Today.Year, 8, 1);
            var startDate = endDate.AddYears(-14);
            if (cachedUnEmpData == null || cachedUnEmpData.Count() <= 0)
            {

                var unEmpData = ParseData.CallQuandlAPI(startDate, endDate);

                cachedUnEmpData = unEmpData;
            }


            var data = cachedUnEmpData.Where(x => x.Date >= date.AddDays(-2) && x.Date <= date.AddDays(2)).Select(x => new { x.State, x.Value })
                .ToDictionary(x => x.State, x => (object)((int)x.Value));
            //var modData = new List<KeyValuePair<string, object>>() { (new KeyValuePair<string, object>("State", "Unemployment")) }.Union(data).ToDictionary();
            ;


            return Json(ConvertObjToSerializedArr(data));
        }

        [Route("getStateInfo")]
        public IHttpActionResult GetStateInfo([FromUri] string state)
        {
            return Json(state);
            //var endDate = new DateTime(DateTime.Today.Year, DateTime.Today.Month, 1);
            //var startDate = endDate.AddYears(-2);
            //if (cachedUnEmpData == null || cachedUnEmpData.Count() <= 0)
            //{

            //    var unEmpData = ParseData.CallQuandlAPI(startDate, endDate);

            //    cachedUnEmpData = unEmpData;
            //}


            //var data = cachedUnEmpData.Where(x => x.Date == date).Select(x => new { x.State, x.Value })
            //    .ToDictionary(x => x.State, x => (object)((int)x.Value));
            ////var modData = new List<KeyValuePair<string, object>>() { (new KeyValuePair<string, object>("State", "Unemployment")) }.Union(data).ToDictionary();
            //;


            //return Json(ConvertObjToSerializedArr(data));


        }

        private List<object[]> ConvertObjToSerializedArr(Dictionary<string, object> data)
        {
            List<object[]> dataSet = new List<object[]>();
            foreach (var dataPair in data)
            {
                var itemElement = new object[2] { dataPair.Key, dataPair.Value };
                dataSet.Add(itemElement);
            }

            dataSet.Insert(0, new object[] { "State", "Unemployment" });

            return dataSet;
        }

        

    }
}
